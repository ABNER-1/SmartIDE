package i18n

import (
	"embed"

	"encoding/json"

	"strings"

	"github.com/jeandeaual/go-locale"

	// 这里不要忘记引入驱动,如引入默认的json驱动

	_ "github.com/leansoftX/i18n/parser_json"
)

//
type I18nSource struct {
	Config struct {
		Info_help_short         string `json:"info_help_short"`
		Info_help_long          string `json:"info_help_long"`
		Info_set_config_success string `json:"info_set_config_success"`
		Err_read_config         string `json:"err_read_config"`
		Err_set_config          string `json:"err_set_config"`

		Info_read_docker_compose      string `json:"info_read_docker_compose"`
		Err_services_not_exit         string `json:"err_services_not_exit"`
		Err_file_not_exit             string `json:"err_file_not_exit"`
		Err_devcontainer_not_contains string `json:"err_devcontainer_not_contains"`
		Err_Gitconfig_not_exit        string `json:"err_gitconfig_not_exit"`
	} `json:"config"`

	Main struct {
		Info_Info_help_short      string `json:"info_help_short"`
		Info_Info_help_long       string `json:"info_help_long"`
		Info_Info_help_flag_debug string `json:"info_help_flag_debug"`
		Info_Usage_template       string `json:"info_usage_template"`
		Info_workspace_loading    string `json:"info_workspace_loading"`
		Info_ssh_connect_check    string `json:"info_ssh_connect_check"`

		Err_file_not_exit       string `json:"err_file_not_exit"`
		Err_file_not_exit2      string `json:"err_file_not_exit2"`
		Err_version_not_build   string `json:"err_version_not_build"`
		Err_flag_value_invalid  string `json:"err_flag_value_invalid"`
		Err_flag_value_invalid2 string `json:"err_flag_value_invalid2"`
		Err_flag_value_required string `json:"err_flag_value_required"`
		Err_workspace_none      string `json:"err_workspace_none"`
	} `json:"main"`

	Help struct {
		Info_help_short string `json:"info_help_short"`
		Info_help_long  string `json:"info_help_long"`
	} `json:"help"`

	VmStart struct {
		Info_starting              string `json:"info_starting"`
		Info_connect_remote        string `json:"info_connect_remote"`
		Info_git_clone             string `json:"info_git_clone"`
		Info_git_checkout_and_pull string `json:"info_git_checkout_and_pull"`
		Info_read_config           string `json:"info_read_config"`
		Info_create_network        string `json:"info_create_network"`
		Info_compose_up            string `json:"info_compose_up"`
		Info_warting_for_webide    string `json:"info_warting_for_webide"`
		Info_open_brower           string `json:"info_open_brower"`
		Info_git_cloned            string `json:"info_git_cloned"`
	} `json:"vm_start"`

	Version struct {
		Info_help_short string `json:"info_help_short"`
		Info_help_long  string `json:"info_help_long"`
		Info_template   string `json:"info_template"`
	} `json:"version"`

	Update struct {
		Info_help_short        string `json:"info_help_short"`
		Info_help_long         string `json:"info_help_long"`
		Info_remove_repeat     string `json:"info_remove_repeat"`
		Info_help_flag_version string `json:"info_help_flag_version"`
		Info_help_flag_build   string `json:"info_help_flag_build"`

		Warn_rel_lastest string `json:"warn_rel_lastest"`
	} `json:"update"`

	Start struct {
		Info_help_short         string `json:"info_help_short"`
		Info_help_long          string `json:"info_help_long"`
		Info_help_flag_host     string `json:"info_help_flag_host"`
		Info_help_flag_port     string `json:"info_help_flag_port"`
		Info_help_flag_username string `json:"info_help_flag_username"`
		Info_help_flag_password string `json:"info_help_flag_pasword"`
		Info_help_flag_repourl  string `json:"info_help_flag_repourl"`
		Info_help_flag_branch   string `json:"info_help_flag_branch"`
		Info_help_flag_filepath string `json:"info_help_flag_filepath"`

		Info_start                   string `json:"info_start"`
		Info_end                     string `json:"info_end"`
		Info_running_container       string `json:"info_running_container"`
		Info_running_openbrower      string `json:"info_running_openbrower"`
		Info_docker_compose_filepath string `json:"info_docker_compose_filepath"`
		Info_ssh_tunnel              string `json:"info_ssh_tunnel"`
		Info_create_network          string `json:"info_create_network"`
		Info_open_in_brower          string `json:"info_open_in_brower"`
		Info_workspace_saving        string `json:"info_workspace_saving"`
		Info_workspace_saved         string `json:"info_workspace_saved"`

		Err_docker                    string `json:"err_docker"`
		Err_DockerPs                  string `json:"err_dockerps"`
		Err_Docker_Compose            string `json:"err_docker_compose"`
		Warn_docker_container_started string `json:"warn_docker_container_started"`
		Err_Docker_compose_save       string `json:"err_docker_compose_save"`
	} `json:"start"`

	List struct {
		Info_help_short string `json:"info_help_short"`
		Info_help_long  string `json:"info_help_long"`
		Info_start      string `json:"Info_start"`
		Info_end        string `json:"info_end"`

		Info_dal_none              string `json:"info_dal_none"`
		Info_workspace_list_header string `json:"info_workspace_list_header"`
	} `json:"list"`

	Get struct {
		Info_help_short            string `json:"info_help_short"`
		Info_help_long             string `json:"info_help_long"`
		Info_help_flag_workspaceid string `json:"info_help_flag_workspaceid"`

		Info_workspace_detail_template      string `json:"info_workspace_detail_template"`
		Info_workspace_host_detail_template string `json:"info_workspace_host_detail_template"`

		Warn_flag_workspaceid_none string `json:"warn_flag_workspaceid_none"`
	} `json:"get"`

	Stop struct {
		Info_help_short         string `json:"info_help_short"`
		Info_help_long          string `json:"info_help_long"`
		Info_start              string `json:"info_start"`
		Info_end                string `json:"info_end"`
		Info_help_flag_filepath string `json:"info_help_flag_filepath"`

		Info_sshremote_connection_creating string `json:"info_sshremote_connection_creating"`
		Info_docker_stopping               string `json:"info_docker_stopping"`
	} `json:"stop"`

	Remove struct {
		Info_help_short       string `json:"info_help_short"`
		Info_help_long        string `json:"info_help_long"`
		Info_start            string `json:"info_start"`
		Info_end              string `json:"info_end"`
		Info_flag_workspaceid string `json:"info_flag_workspaceid"`
		Info_flag_yes         string `json:"info_flag_yes"`
		Info_flag_force       string `json:"info_flag_force"`

		Info_sshremote_connection_creating string `json:"info_sshremote_connection_creating"`
		Info_docker_removing               string `json:"info_docker_removing"`
		Info_project_dir_removing          string `json:"info_project_dir_removing"`

		Info_is_confirm_remove string `json:"info_is_confirm_remove"`

		Err_workspace_not_exit string `json:"err_workspace_not_exit"`
	} `json:"remove"`

	New struct {
		Info_help_short              string `json:"info_help_short"`
		Info_help_long               string `json:"info_help_long"`
		Info_help_info               string `json:"info_help_info"`
		Info_help_info_operation     string `json:"info_help_info_operation"`
		Info_help_flag_type          string `json:"info_help_flag_type"`
		Info_help_flag_projectFolder string `json:"info_help_flag_projectFolder"`
		Info_yaml_exist              string `json:"info_yaml_exist"`
		Info_noempty_is_comfirm      string `json:"info_noempty_is_comfirm"`
		Info_type_no_exist           string `json:"info_type_no_exist"`
		Info_creating_project        string `json:"info_creating_project"`
		Info_loading_templates       string `json:"info_loading_templates"`
		Info_templates_list_header   string `json:"info_templates_list_header"`
		Err_read_templates           string `json:"err_read_templates"`
	} `json:"new"`

	Host struct {
		Info_help_short       string `json:"info_help_short"`
		Info_help_long        string `json:"info_help_long"`
		Info_help_get_short   string `json:"info_help_get_short"`
		Info_help_get_long    string `json:"info_help_get_long"`
		Info_help_list_short  string `json:"info_help_list_short"`
		Info_help_list_long   string `json:"info_help_list_long"`
		Info_help_flag_hostid string `json:"info_help_flag_hostid"`

		Info_host_table_header    string `json:"info_host_table_header"`
		Info_host_detail_template string `json:"info_host_detail_template"`

		Err_host_data_not_exit string `json:"err_host_data_not_exit"`
	} `json:"host"`

	Common struct {
		Debug_key_public           string `json:"debug_key_public"`
		Debug_same_not_overwrite   string `json:"debug_same_not_overwrite"`
		Debug_auto_connect_gitrepo string `json:"debug_auto_connect_gitrepo"`
		Debug_empty_error          string `json:"debug_empty_error"`

		Err_sshremote_param_repourl_none string `json:"err_sshremote_param_repourl_none"`
		Err_password_none                string `json:"err_password_none"`
		Err_dal_record_repeat            string `json:"err_dal_record_repeat"`
		Err_dal_update_fail              string `json:"err_dal_update_fail"`
		Err_dal_update_count_too_much    string `json:"err_dal_update_count_too_much"`
		Err_enum_error                   string `json:"err_ernum_error"`
		Err_ssh_password_required        string `json:"err_ssh_password_required"`
		Err_git_env_check                string `json:"err_git_env_check"`

		Info_privatekey_is_overwrite      string `json:"info_privatekey_is_overwrite"`
		Info_whether_overwrite            string `json:"info_whether_overwrite"`
		Info_gitrepo_clone_done           string `json:"info_gitrepo_clone_done"`
		Info_gitrepo_cloned               string `json:"info_gitrepo_cloned"`
		Info_please_enter_password        string `json:"info_please_enter_password"`
		Info_canel_privatekey_password    string `json:"info_canel_privatekey_password"`
		Info_port_is_binding              string `json:"info_port_is_binding"`
		Info_port_binding_result2         string `json:"info_port_binding_result2"`
		Info_port_binding_result          string `json:"info_port_binding_result"`
		Info_find_new_port                string `json:"info_find_new_port"`
		Info_ssh_webide_host_port         string `json:"info_ssh_webide_host_port"`
		Info_ssh_host_port                string `json:"info_ssh_host_port"`
		Info_temp_create_directory        string `json:"info_temp_create_directory"`
		Info_temp_created_docker_compose  string `json:"info_temp_created_docker_compose"`
		Info_temp_created_config          string `json:"info_temp_created_config"`
		Info_table_header_containers      string `json:"info_table_header_containers"`
		Info_ssh_rsa_cancel_pwd_successed string `json:"info_ssh_rsa_cancel_pwd_successed"`

		Warn_dal_record_not_exit_condition string `json:"warn_dal_record_not_exit_condition"` // 没有查询到对应的数据
		Warn_dal_record_not_exit           string `json:"warn_dal_record_not_exit"`           // 没有查询到对应的数据
		Warn_param_is_null                 string `json:"warn_param_is_null"`                 // 参数为空

	} `json:"common"`
}

var instance *I18nSource

/* var I18nSource_EN string
var I18nSource_ZH string */

//go:embed language/*
var f embed.FS

// get internationalization source
// 获取当前系统的语言，动态加载对应的json文件并解析成结构体，方便在代码中调用
// 1. 新增，首先在 “lib/i18n/language” 的对应节点下新增，并同步在 “lib/i18n/language/translate.go” 中的 “I18nSource” 增加相应的属性；
// 2. 在代码中使用
//    var instanceI18nStart = i18n.GetInstance().Start
//    fmt.println(instanceI18nStart.Info_help_short)
func GetInstance() *I18nSource {
	if instance == nil {
		// locale
		currentLang, _ := locale.GetLocale()
		if strings.Index(currentLang, "zh-") == 0 { // 如果不是简体中文，就是英文
			currentLang = "zh_cn"
		} else {
			currentLang = "en_us"
		}

		// loading and parse json
		data, _ := f.ReadFile("language/" + currentLang + "/info.json")
		json.Unmarshal(data, &instance)
	}
	return instance
}
